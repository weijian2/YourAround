package com.around;

import com.google.api.services.bigquery.model.TableFieldSchema;
import com.google.api.services.bigquery.model.TableRow;
import com.google.api.services.bigquery.model.TableSchema;
import com.google.cloud.bigtable.dataflow.CloudBigtableIO;
import com.google.cloud.bigtable.dataflow.CloudBigtableScanConfiguration;
import com.google.cloud.dataflow.sdk.Pipeline;
import com.google.cloud.dataflow.sdk.io.BigQueryIO;
import com.google.cloud.dataflow.sdk.io.Read;
import com.google.cloud.dataflow.sdk.options.PipelineOptions;
import com.google.cloud.dataflow.sdk.options.PipelineOptionsFactory;
import com.google.cloud.dataflow.sdk.transforms.DoFn;
import com.google.cloud.dataflow.sdk.transforms.ParDo;
import com.google.cloud.dataflow.sdk.values.PCollection;
import org.apache.hadoop.hbase.client.Result;
import org.apache.hadoop.hbase.util.Bytes;

import java.nio.charset.Charset;
import java.util.ArrayList;
import java.util.List;

public class PostDumpFlow {
    private static final String PROJECT_ID = "youraround-cmu";
    private static final Charset UTF8_CHARSET = Charset.forName("UTF-8");

    public static void main(String[] args) {
        // reading from Cloud Bigtable: https://cloud.google.com/bigtable/docs/dataflow-hbase#reading
        CloudBigtableScanConfiguration config = new CloudBigtableScanConfiguration.Builder()
                .withProjectId(PROJECT_ID)
                .withInstanceId("around-post")
                .withTableId("post")
                .build();

        // Start by defining the options for the pipeline.
        PipelineOptions options = PipelineOptionsFactory.fromArgs(args).create();
        Pipeline p = Pipeline.create(options);
        // E(Extracted), read bigTable data row by row, Result is byte array
        PCollection<Result> btRows = p.apply(Read.from(CloudBigtableIO.read(config)));
        // T(Transformed), transform BigTable row(PCollection<Result>) to BigQuery row(PCollection<TableRow>), TableRow is row of BigQuery
        PCollection<TableRow> bqRows = btRows.apply(ParDo.of(new DoFn<Result, TableRow>() {
            @Override
            public void processElement(ProcessContext c) {
                Result result = c.element();
                String postId = new String(result.getRow());
                String user = new String(result.getValue(Bytes.toBytes("post"), Bytes.toBytes("user")), UTF8_CHARSET);
                String message = new String(result.getValue(Bytes.toBytes("post"), Bytes.toBytes("message")), UTF8_CHARSET);
                String lat = new String(result.getValue(Bytes.toBytes("location"), Bytes.toBytes("lat")), UTF8_CHARSET);
                String lon = new String(result.getValue(Bytes.toBytes("location"), Bytes.toBytes("lon")), UTF8_CHARSET);
                //BQ Table row object
                TableRow row = new TableRow();
                row.set("postId", postId);
                row.set("user", user);
                row.set("message", message);
                row.set("lat", Double.parseDouble(lat));
                row.set("lon", Double.parseDouble(lon));
                c.output(row);
            }
        }));
        // L(load), write data into BigQuery: https://cloud.google.com/dataflow/model/bigquery-io#writing-to-bigquery
        List<TableFieldSchema> fields = new ArrayList<>();
        fields.add(new TableFieldSchema().setName("postId").setType("STRING"));
        fields.add(new TableFieldSchema().setName("user").setType("STRING"));
        fields.add(new TableFieldSchema().setName("message").setType("STRING"));
        fields.add(new TableFieldSchema().setName("lat").setType("FLOAT"));
        fields.add(new TableFieldSchema().setName("lon").setType("FLOAT"));

        TableSchema schema = new TableSchema().setFields(fields);
        bqRows.apply(BigQueryIO.Write
                .named("Write")
                .to(PROJECT_ID + ":" + "post_analysis" + "." + "daily_dump_" + System.currentTimeMillis())
                .withSchema(schema)
                .withWriteDisposition(BigQueryIO.Write.WriteDisposition.WRITE_TRUNCATE)
                .withCreateDisposition(BigQueryIO.Write.CreateDisposition.CREATE_IF_NEEDED));
        // run pipeline
        p.run();
    }
}

# Introduction
YourAround is a geo-index based Social Network based on Google Cloud Platform. Front-end is implemented by React.js, back-end is implemented by Golang. User can post posts to my platform, see nearby posts by other users in Google Map. User can filter some posts they don't want to see by setting keywords themselves. User can also search all posts in specific position with specified distance range.

## Why this project
Social Networking Services(SNS), such as WeChat and Facebook, have been trending around the world for years. Now the industry is looking for the next generation of SNS and its applications. In my opinion, the opportunity lies in improved user experience(like filter fake news and spam ads) and increased content relevancy(timeline based showing has been gradually eliminated). This two pain points has not been solved by WeChat and Facebook, so I want to make a SNS application that show contents based on geo-index and can support spam filter. This project contains all of my understanding of a social network.

## Tech Stack
* Go/React.js
* ElasticSearch/BigTable/DataFlow/BigQuery
* Google App Engine(GAE)/Google Compute Engine(GCE)/Google Cloud Storage(GCS)

## Project structure
The overall structure (tech stack) of this project:
![](https://github.com/weijian2/YourAround/raw/master/demoPics/structure.png)
![](https://github.com/weijian2/YourAround/raw/master/demoPics/dataflow.png)

Tech Stacks choosen reasons<br>
**1. Why choose Golang as backend language?**<br>
* Go is built for system, web app, pretty low level, we don’t need Django for python and TomCat for java. It has built-in web services to handle HTTP request/response, URL mapping…
* Go has been very well designed and optimised for scaling(Go routine and channel)
* Go is a fast, statically typed, compiled language that feels like a dynamically typed, interpreted language.
* In my opinion, go the future of server language.

**2. Why choose ElasticSearch?**<br>
* ES is an open source, distributed, RESTful search engine. As the heart of the Elastic Stack, it centrally stores our data so we can query the data quickly.
* ES is implemented by K-D tree, which is very efficient to perform 2-dimensional search problem(e.g. given latitute and longitute, search all posts in that position within 200km)
* I also use ES to store structured data(post's content, lat and lon)

**3. Why choose Google App Engine Flex?**<br>
* Google App Engine is a cloud computing platform for developing and hosting web applications in Google-managed data centers.
* Actually, GAE is based on GCE, however, by using GAE, I don't need to build a system by myself, load balancer and VMs are implemented in black box in GAE.
* App Engine flexible environment automatically scales app up and down while balancing the load. Microservices, authorization, SQL and NoSQL databases, traffic splitting, logging. 

**4. Why choose Google Compute Engine?**<br>
* GCE is similar to Amazon EC2, which offers virtual machines (Xen, KVM, VirtualBox, VMware ESX, etc.). Amazon EC2 and Google Compute Engine belong to IaaS as well. 
* I install ES in the VM(Ubuntu 16) in the GCE to store post information and perform search.

**5. Why use Google Cloud Storage(GCS)?**<br>
* To store unstructured data (user posted images). BigTable, BigQuery, ElasticSearch are all for structured data. GCS is well-known for its durability, scalability and availability.

**6. Why use BigTable instead of MySQL or MangoDB?**<br>
* I want to make this project scalable, MySQL do not have enough scability. MongoDB is good, but it is just software, in order to make it scalable, I need to rent vms and clusters, to build system by myself. BigTable is a sparse, distributed, persistent multidimensional sorted map. It runs on Google cloud, which is very scalable, I only need to implement my business logic, which is very convenient. 

**7. Why use BigTable since I already use ES to save post data?**<br>
* I want to do offline data analysis on it using BigQuery. ES is search engine with complicated query support and better read performance but may lose data, so I save posts data on ES and Bigtable both.
* Save post data to ES is for geo-index based range search. Save post data to BigTable is for offline data analysis using BigQuery latter.

**8. Why use DataFlow and BigQuery?**<br>
* I want to do some data analysis of user posted data(like sentiment analysis, spam filter..). BigTable is designed for saving large volume of data but not good for querying (JOIN, ORDER BY etc.). BigQuery, instead, is efficient for querying. It is Google's fully managed, petabyte scale, low cost enterprise data warehouse. It supports data analysis by using SQL syntax.
* In order to dump data into BigQuery from BigTable, I need to use Google DataFlow, which runs in Google's cloud platform, who can help me to transform my data from BigTable to BigQuery easily.

**9. Why use Cache?**<br>
* When performance needs to be improved, caching is often the first step to take, so I decide to place previously requested information (search query) in cache to speed up data access and reduce bandwidth demand. In this project, I use TTL(time to live) as caching Strategy.

**10. Why Choose Redis instead of Memcached?**<br>
* Redis is an open source, in-memory data structure store, used as a database, cache and message broker. Basically a key-values store. Google Cloud has a memcache similarly but it does not support GAE flex yet.
* Redis has more features than Memcached and is, thus, more powerful and flexible. Refered this link: https://www.infoworld.com/article/3063161/nosql/why-redis-beats-memcached-for-caching.html

## Tools
GoLand
Google Cloud Platform
Apache Maven(it will install all dependencies(dataflow related libraries) automatically, and can create Java project and manage it for you)

## Requirements

## Installation

## Usage/Quick Start

## Screenshots

## Known bugs
NAN

## Todo list
* Front-end implementation using React.js.

## Deployment
Deployed backend to GAE(Google App Engine)<br>
click [demo link](https://youraround-cmu.appspot.com) to visit

## Change Log
1. v0.0.1(02/11/2018)<br>
* Create GCE(Google Compute Engine) and install elastic search(golang version) in it.
* Create handlerSearch and handlerPost functions.
2. v0.0.2(02/12/2018)<br>
* Implement handlerSearch, handlerPost functions.
3. v0.0.3(02/15/2018)<br>
* Add filter sensitive words function in search method(Spam and Abuse Detection).
* Deploy backend to Google App Engine.
4. v0.0.4(02/21/2018)<br>
* Create bucket in Google Cloud Storage(GCS).
* Update handlerPost function to parse multipart form and store post with image file into GCS.
5. v0.0.5(02/23/2018)<br>
* Update handlerPost function to support save post data into BigTable.
6. v0.0.6(02/25/2018)<br>
* Implement Google dataflow to dump post data from BigTable to BigQuery.
7. v0.0.7(03/01/2018)<br>
* Used OAuth 2.0 to support authentication and user signup.
8. v0.0.8(03/10/2018)<br>
* Used Redis to cache search requests to speed up user query.
9. v0.0.9(03/17/2018)<br>
* Init front-end React project, front-end change log will updated in /youraround-front.

## Licenses
NAN

## Notes
If you have any questions, feel free to contact me at weijian1@andrew.cmu.edu


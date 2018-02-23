# Introduction


## Why this project


## Tech Stack
* Go/React.js
* BigTable/BigQuery/ElasticSearch
* Google App Engine(GAE)/Google Compute Engine(GCE)/Google Cloud Storage(GCS)

## Project structure
The overall structure (tech stack) of this project:
![](https://github.com/weijian2/YourAround/raw/master/demoPics/structure.png)

Tech Stack choosen<br>
**1. Why choose Golang as backend language?**<br>
* Go is built for system, web app, pretty low level, we don’t need Django for python and TomCat for java. It has built-in web services to handle HTTP request/response, URL mapping…
* Go has been very well designed and optimised for scaling(Go routine and channel)
* Go is a fast, statically typed, compiled language that feels like a dynamically typed, interpreted language.
* In my opinion, go the future of server language.

**2. Why choose elasticSearch?**<br>
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

**5. Why Google Cloud Storage(GCS)?**<br>
* To store unstructured data (user posted images). BigTable, BigQuery, ElasticSearch are all for structured data. GCS is well-known for its durability, scalability and availability.

**6. Why use BigTable?**<br>
* BigTable is a sparse, distributed, persistent multidimensional sorted map. We want to do offline data analysis on it using BigQuery. ES is search engine with complicated query support and better read performance but may lose data, so I store posts data on ES and Bigtable both.
* Save post data to ES is for geo-index based range search. Save post data to BigTable is for offline data analysis using BigQuery.

**7. Why use BigQuery?**<br>

## Requirements

## Installation

## Usage/Quick Start

## Screenshots

## Known bugs
NAN.
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

## Licenses
NAN

## Notes
If you have any questions, feel free to contact me at weijian1@andrew.cmu.edu


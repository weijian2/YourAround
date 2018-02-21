# Introduction


## Why doing this project


## Tech Stack
* Go/React.js
* BigTable/BigQuery/ElasticSearch
* Google App Engine(GAE)/Google Compute Engine(GCE)/Google Cloud Storage(GCS)

## Project structure
The overall structure (tech stack) of this project:
![](https://github.com/weijian2/YourAround/raw/master/demoPics/structure.png)
why Google Cloud Storage(GCS)?<br>
* To store unstructured data (user posted images). BigTable, BigQuery, ElasticSearch are all for structured data. GCS is well-known for its durability, scalability and availability.
why ElasticSearch?<br>
* ...

## Requirements

## Installation

## Usage/Quick Start

## Screenshots

## Known bugs

## Todo list

## Deployment
Deployed backend to GAE(Google App Engine)<br>
click [demo link](https://youraround-cmu.appspot.com) to visit

## Change Log
1. v0.0.1(02/11/2018)<br>
* Create GCE(Google Compute Engine) and install elastic search(golang version) in it
* Create handlerSearch and handlerPost functions
2. v0.0.2(02/12/2018)<br>
* Implement handlerSearch, handlerPost functions
3. v0.0.3(02/15/2018)<br>
* Add filter sensitive words function in search method(Spam and Abuse Detection)
* Deploy backend to Google App Engine
4. v0.0.4(02/21/2018)<br>
* Create bucket in Google Cloud Storage(GCS)
* Update handlerPost function to parse multipart form and store post with image file into GCS

## Licenses
NAN

## Notes



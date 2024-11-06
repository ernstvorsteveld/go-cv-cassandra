docker run -it --rm --name logstash --net docker_default opensearchproject/logstash-oss-with-opensearch-output-plugin:7.16.2 -e 'input { stdin { } } output {
   opensearch {
     hosts => ["https://opensearch-node1:9200"]
     index => "opensearch-logstash-docker-%{+YYYY.MM.dd}"
     user => "admin"
     password => "Pass*w0rd!"
     ssl => true
     ssl_certificate_verification => false
   }
 }'




 docker run --rm -it -v ~/pipeline/:/usr/share/logstash/pipeline/ docker.elastic.co/logstash/logstash:8.15.

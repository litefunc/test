# docker run -dit --name tecmint-web -p 8080:80 -v (pwd)/htdocs/:/usr/local/apache2/htdocs/ httpd:2.4
docker run -d --name tecmint-web -p 8080:80 -v (pwd)/htdocs/:/usr/local/apache2/htdocs/ httpd:2.4

docker run -dit --name my-apache-app -p 8080:80  aarch64/httpd:2.4

docker run -dit --name my-apache-app -p 8086:80  aarch64/httpd:2.4

docker run -d --name my-apache-app --network host  aarch64/httpd:2.4


docker run -d --name my-apache-app --network host -v (pwd)/htdocs/:/usr/local/apache2/htdocs/ apache:2.4
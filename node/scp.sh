scp Hub/Dockerfile2 root@192.168.2.2:~/Hub/Dockerfile2
scp Hub/Dockerfile3 root@192.168.2.2:~/Hub/Dockerfile3
scp Hub/resolv.conf root@192.168.2.2:~/Hub/resolv.conf

docker run -it -v /root/Hub/resolv.conf:/etc/resolv.conf --name web -p 8888:8888 web:latest sh

scp ~/program/temp/ms/web1.tar root@192.168.6.141:~/web1.tar

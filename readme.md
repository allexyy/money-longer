docker run -d \
--name monyLonger-db \
-e POSTGRES_USER=admin \
-e POSTGRES_PASSWORD=secret \
-e POSTGRES_DB=monyLonger \
-p 5432:5432 \
postgres:16-alpine


docker exec -it monyLonger-db psql -U admin -d monyLonger

docker stop monyLonger-db
docker start monyLonger-db
for i in {1..5};
do
docker compose -f docker-compose.test.yml down -v
make test
done

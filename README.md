# Master git

Be good at dealing with CSVs

Run docker for RabbitMQ:

```BASH
docker run -d --hostname my-rabbit --name some-rabbit -p 4369:4369 -p 5671:5671 -p 5672:5672 -p 15672:15672 rabbitmq
docker exec some-rabbit rabbitmq-plugins enable rabbitmq_management
```

## Working with Makefile 

Commit your changes before building new release!  

```BASH
make build
```

does all necessary things.
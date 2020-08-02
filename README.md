 # copper-face-jacks

`copper-face-jacks` serves as an implementation of a Kafka consumer application.

## Goals

 1. Subscribes to `test_topic` topic and consumes messages
 2. Parses the message
 3. Sends message to redis

To run the full capability, you'll need both the producer [hairy-lemon](https://github.com/lhmzhou/hairy-lemon) and the consumer [copper-face-jacks](https://github.com/lhmzhou/copper-face-jacks). `copper-face-jacks` is only 50% of the capability.


## System Architecture

Below is the bird's eye view of the e2e initiative. The red enclosure represents `copper-face-jacks`.

![alt diagram](https://github.com/lhmzhou/copper-face-jacks/blob/master/image/cfj_arch.png)


## Prerequisites

- go v.1.13.1+
- kafka container from [confluentinc's docker-compose.yaml](https://github.com/confluentinc/examples/blob/5.3.1-post/cp-all-in-one/docker-compose.yml)
- [dep](https://github.com/golang/dep) for managing dependencies
- [redis container](https://hub.docker.com/_/redis)

## Build

1. Initialize redis instance by running below commands in cli:

    ```
    $ docker pull redis
    $ docker run --name testredis -p 6379:6379 -d redis
    ```

2. For initializing kafka instance, copy the [yaml file](https://github.com/confluentinc/examples/blob/5.3.1-post/cp-all-in-one/docker-compose.yml) to desktop.

    In the yaml file, delete everything apart from `zookeeper` and `broker` configs. Save the file. Next, execute below commands in cli:

    ```
    $ docker-compose up
    $ docker ps
    ```

3. Clone the [producer](https://github.com/lhmzhou/hairy-lemon).

    Run `hairy-lemon`. Open browser and hit the api at `http://localhost:8081/score/1`

4. Clone the [consumer](https://github.com/lhmzhou/copper-face-jacks).

    Run `copper-face-jacks`, which will consume the message from kafka producer, parse it, and send it to redis.

### Test data in redis

Go to redis docker instance by executing below command in cli:

```
$ docker exec -it testredis /bin/bash
$ redis-cli
$ keys *
```

From `key*` command above, get list of keys that are sent to redis. In our case, you will see key as 1. To get all the values associated with key 1, run: 

```
$ get 1
```

## h/t

Maximum respect and many thanks to the developers on these open-source projects for making `copper-face-jacks` possible:

[rcrowley/go-metrics](https://github.com/rcrowley/go-metrics)
</br>
[jcmturner/gofork](https://github.com/jcmturner/gofork)
</br>
[klauspost/compress](https://github.com/klauspost/compress)
</br>
[go-redis/redis](https://github.com/go-redis/redis)
</br>
[golang/snappy](https://github.com/golang/snappy)
</br>
[hashicorp/go-uuid](https://github.com/hashicorp/go-uuid)
</br>
[pierrec/lz4](https://github.com/pierrec/lz4)
</br>
[eapache @ Shopify](https://github.com/Shopify/sarama)
</br>
[davecgh/go-spew](https://github.com/davecgh/go-spew)

# Redis at home

![for_exceptional_broski](https://github.com/theredditbandit/redis-at-home/assets/85390033/17cec56f-1ae6-46c0-b19e-99222e8ec942)



# Features

- Full support for the [PING](https://redis.io/commands/ping/) command
- Full support for the [ECHO](https://redis.io/commands/echo/) command
- partial support for the [SET](https://redis.io/commands/set/) command
    - PX flag is supported
- Full support for the [GET](https://redis.io/commands/get/) command


# TODO
 - Write tests for existing commands
 - add the following commands
    -EXISTS - check if a key is present.
    -DEL - delete one or more keys.
    -INCR - increment a stored number by one.
    -DECR - decrement a stored number by one.
    -LPUSH - insert all the values at the head of a list.
    -RPUSH - insert all the values at the tail of a list.
    -SAVE - save the database state to disk, you should also implement load on startup alongside this.
    -FLUSHDB - flush the database.
    -HSET - set a field to a value
 - add video recording to readme
 - upload stats from redis-benchmark comparing this implementation with actual redis

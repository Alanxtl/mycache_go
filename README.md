# My cache

A go implemented distributed cache follow the tutorial of [geektutu](https://geektutu.com/post/geecache.html).

The new features I've added on top of geektutu:

1. Using Nacos as the register center of distributed cache.
2. Multiple ways of peer communicate, http and dubbo. 

### Prerequisites

1. Clone my project

    ```bash
    git clone https://github.com/Alanxtl/mycache_go.git
    cd mycache_go
    ```

2. Docker

    ```bash
    docker-compose up
    ```

3. Start test

    ```bash
    ./run.sh
    ```

### Credit

1. https://geektutu.com/post/geecache.html
2. https://github.com/1055373165/ggcache
3. https://github.com/peanutzhen/peanutcache

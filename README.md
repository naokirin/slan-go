# slan-go

slan-go is a slack bot.

## Usage

1. Add `config/slan-go.conf`
2. Build docker image and run container

### Prepare slan-go.conf

Write slack bot token, mention name and other settings in `config/slan-go.conf` .

See `config/slan-go.conf.sample` .

TODO: add how to write.

### Build docker image and Run container

```
docker build -t slan-go .
docker run --rm --name slan-go slan-go
```

## License

See https://github.com/naokirin/slan-go/blob/master/LICENSE


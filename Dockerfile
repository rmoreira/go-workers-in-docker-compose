FROM wernight/phantomjs
COPY vendor/github.com/nladuo/go-phantomjs-fetcher/phantomjs_fetcher.js phantomjs_fetcher.js
COPY main main
CMD ["./main"]

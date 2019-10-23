FROM arm32v7/ruby:2.6.3-alpine

COPY . /app
WORKDIR /app

ENV RAILS_ENV production

RUN apk update && \
    apk upgrade && \
    apk add build-base sqlite-dev tzdata && \
    rm -rf /var/cache/apk/*

RUN bundle install --deployment --without development test

EXPOSE 3000
CMD bundle exec rails server -b 0.0.0.0
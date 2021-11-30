# Example project with Litestream and Go

This repo contains a simple application that tracks number of page views. The main purpose was to prototype running [Litestream](https://litestream.io/), replicating a local (to the application) SQLite database to persistent storage in Google Cloud Storage. 

### Deployment

To deploy to Google Cloud

```
PROJECT=$(gcloud config get-value project)
NAME=litestream-demo
TAG=gcr.io/$PROJECT/$NAME
gcloud builds submit --tag $TAG
gcloud beta run deploy $NAME --image $TAG \
            --platform=managed \
            --region=europe-west1 \
            --execution-environment gen2
```

#### Environment variables

The application needs the following environment variables.

```
REPLICA_URL=gcs://go-litestream
DB_PATH=db.sqlite
LITESTREAM_ACCESS_KEY_ID=GOOG1EX....
LITESTREAM_SECRET_ACCESS_KEY=VWq....
```


## 環境構築手順

```bash
git clone https://github.com/Tsutomu-Ikeda/sake-log
cd sake-log
docker-compose build
docker-compose up -d
http --print=b localhost:8080/list
# [
#     {
#         "firstName": "hoge",
#         "id": 1,
#         "lastName": "fuga"
#     },
#     {
#         "firstName": "hoge",
#         "id": 2,
#         "lastName": "fuga"
#     },
#     {
#         "firstName": "hoge",
#         "id": 3,
#         "lastName": "fuga"
#     }
# ]
http --print=b POST localhost:8080/insert firstname=hoge lastname=fuga
# New user was added
http --print=b localhost:8080/list                          
# [
#     {
#         "firstName": "hoge",
#         "id": 1,
#         "lastName": "fuga"
#     },
#     {
#         "firstName": "hoge",
#         "id": 2,
#         "lastName": "fuga"
#     },
#     {
#         "firstName": "hoge",
#         "id": 3,
#         "lastName": "fuga"
#     },
#     {
#         "firstName": "hoge",
#         "id": 4,
#         "lastName": "fuga"
#     }
# ]
```

## GCP上への構築

```bash
export PROJECT_ID="fill-here"
export BUCKET_NAME="test-$(date +"%Y%m%d_%H%M%S")"

gcloud storage buckets create gs://$BUCKET_NAME --location asia-northeast1
gcloud iam service-accounts create fs-identity
gcloud projects add-iam-policy-binding $PROJECT_ID --member "serviceAccount:fs-identity@$PROJECT_ID.iam.gserviceaccount.com" --role "roles/storage.objectAdmin"
gcloud alpha run deploy filesystem-app --region asia-northeast1 --source backend/fs-go --execution-environment gen2 --service-account fs-identity --allow-unauthenticated --add-volume=name=gcs,type=cloud-storage,bucket=$BUCKET_NAME --add-volume-mount=volume=gcs,mount-path=/mnt/gcs --update-env-vars MNT_DIR=/mnt/gcs
```

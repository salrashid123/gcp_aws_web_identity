import os
import json
import urllib

import boto3
from botocore.auth import SigV4Auth
from botocore.awsrequest import AWSRequest

def create_token_aws(project_number: str, pool_id: str, provider_id: str) -> None:
    request = AWSRequest(
        method="POST",
        url="https://sts.amazonaws.com/?Action=GetCallerIdentity&Version=2011-06-15",
        headers={
            "Host": "sts.amazonaws.com",
            "x-goog-cloud-target-resource": f"//iam.googleapis.com/projects/{project_number}/locations/global/workloadIdentityPools/{pool_id}/providers/{provider_id}",
        },
    )

    SigV4Auth(boto3.Session().get_credentials(), "sts", "us-east-1").add_auth(request)
    token = {"url": request.url, "method": request.method, "headers": []}
    for key, value in request.headers.items():
        token["headers"].append({"key": key, "value": value})
    #print("%s" % json.dumps(token, indent=2, sort_keys=True))
    print("%s" % urllib.parse.quote(json.dumps(token)))


def main() -> None:

    project_number = os.environ.get('PROJECT_NUMBER')
    pool_id = os.environ.get('POOL_ID')
    provider_id = os.environ.get('PROVIDER_ID')

    create_token_aws(project_number, pool_id, provider_id)


if __name__ == "__main__":
    main()
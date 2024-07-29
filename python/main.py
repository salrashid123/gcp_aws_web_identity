
from google.auth import aws
from google.auth import exceptions
import boto3
from google.cloud import storage

class CustomAwsSecurityCredentialsSupplier(aws.AwsSecurityCredentialsSupplier):
    
    def __init__(self, region):
        self.region = region
 
    def get_aws_security_credentials(self, context, request):
        try:
            session = boto3.Session()
            credentials = session.get_credentials()
            return aws.AwsSecurityCredentials(credentials.access_key, credentials.secret_key, credentials.token)
        except Exception as e:
            raise exceptions.RefreshError(e, retryable=True)

    def get_aws_region(self, context, request):
        return self.region

supplier = CustomAwsSecurityCredentialsSupplier(region='us-east-1')

audience = "//iam.googleapis.com/projects/995081019036/locations/global/workloadIdentityPools/aws-pool-1/providers/aws-provider-1"
subject_token_type = "urn:ietf:params:aws:token-type:aws4_request"
credentials = aws.Credentials(
    audience,
    subject_token_type,
    aws_security_credentials_supplier=supplier
)

# credentials, project = google.auth.default()    
client = storage.Client(project="core-eso",credentials=credentials)
blobs = client.list_blobs("core-eso-bucket")

for blob in blobs:
    print(blob.name)
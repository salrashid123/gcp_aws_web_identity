import * as aw from '@aws-sdk/credential-providers';
import * as gcs from '@google-cloud/storage';
import { AwsClient, AwsSecurityCredentials, AwsSecurityCredentialsSupplier, ExternalAccountSupplierContext } from 'google-auth-library';

let region = 'us-east-1';
let bucketName = 'core-eso-bucket';
let gaudience = "//iam.googleapis.com/projects/995081019036/locations/global/workloadIdentityPools/aws-pool-1/providers/aws-provider-1";

interface MyAwsSupplierOptions {
  credentials?: AwsSecurityCredentials;
  region?: string;
  credentialsError?: Error;
  regionError?: Error;
}

class MyAwsSupplier implements AwsSecurityCredentialsSupplier {
  private readonly credentials?: AwsSecurityCredentials;
  private readonly region?: string;
  private readonly credentialsError?: Error;
  private readonly regionError?: Error;

  constructor(options: MyAwsSupplierOptions) {
    this.credentials = options.credentials;
    this.region = options.region;
    this.credentialsError = options.credentialsError;
    this.regionError = options.regionError;
  }

  async getAwsRegion(context: ExternalAccountSupplierContext): Promise<string> {
    if (this.regionError) {
      throw this.regionError;
    } else {
      return this.region ?? '';
    }
  }

  async getAwsSecurityCredentials(
    context: ExternalAccountSupplierContext
  ): Promise<AwsSecurityCredentials> {
    if (this.credentialsError) {
      throw this.credentialsError;
    } else {
      return this.credentials ?? { accessKeyId: '', secretAccessKey: '' };
    }
  }
}

const provider = aw.fromTokenFile({
  // webIdentityTokenFile: process.env.AWS_WEB_IDENTITY_TOKEN_FILE,
  // roleArn: process.env.AWS_ROLE_ARN,
  // roleSessionName: process.env.AWS_ROLE_SESSION_NAME
});
provider().then((response) => {
  let awsSecurityCredentials = {
    accessKeyId: response.accessKeyId,
    secretAccessKey: response.secretAccessKey,
    token: response.sessionToken
  }

  const clientOptions = {
    audience: gaudience,
    subject_token_type: 'urn:ietf:params:aws:token-type:aws4_request',
    aws_security_credentials_supplier: new MyAwsSupplier({ credentials: awsSecurityCredentials, region: region })
  }

  const authClient = new AwsClient(clientOptions);
  const storage = new gcs.Storage({ authClient });
  async function listFiles() {
    const [files] = await storage.bucket(bucketName).getFiles();
    console.log('Files:');
    files.forEach(file => {
      console.log(file.name);
    });
  }

  listFiles().catch(console.error);

}).catch(console.error);



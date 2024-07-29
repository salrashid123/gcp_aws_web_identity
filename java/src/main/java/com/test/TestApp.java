package com.test;

import com.google.auth.oauth2.ExternalAccountCredentials.SubjectTokenTypes;


import com.google.cloud.storage.Storage;
import com.google.cloud.storage.StorageOptions;

import com.google.api.gax.paging.Page;
import com.google.cloud.storage.Blob;

/*
mvn clean install

 */

public class TestApp {
	public static void main(String[] args) {
		TestApp tc = new TestApp();
	}

	// https://github.com/googleapis/google-auth-library-java/tree/1e27013048d3f73c082d98dd364053da759c5590?tab=readme-ov-file#using-a-custom-supplier-with-aws

	public TestApp() {
		try {


			CustomAwsSupplier awsSupplier = new CustomAwsSupplier();
			com.google.auth.oauth2.AwsCredentials awscredentials = com.google.auth.oauth2.AwsCredentials.newBuilder()
					.setSubjectTokenType(SubjectTokenTypes.AWS4) // Sets the subject token type.
					.setAudience(
							"//iam.googleapis.com/projects/995081019036/locations/global/workloadIdentityPools/aws-pool-1/providers/aws-provider-1") 
					.setAwsSecurityCredentialsSupplier(awsSupplier) // Sets the supplier.
					.build();
;
					Storage storage = StorageOptions.newBuilder().setCredentials(awscredentials).build().getService();
					Page<Blob> blobs = storage.list("core-eso-bucket");
				
					for (Blob blob : blobs.iterateAll()) {
					  System.out.println(blob.getName());
					}

		} catch (Exception ex) {
			System.out.println("Error:  " + ex);
		}
	}

}


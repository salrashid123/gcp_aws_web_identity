package com.test;

import java.io.IOException;

import com.google.auth.oauth2.AwsSecurityCredentials;
import com.google.auth.oauth2.AwsSecurityCredentialsSupplier;
import com.google.auth.oauth2.ExternalAccountSupplierContext;

import software.amazon.awssdk.auth.credentials.AwsCredentials;
import software.amazon.awssdk.auth.credentials.AwsSessionCredentials;
import software.amazon.awssdk.auth.credentials.WebIdentityTokenFileCredentialsProvider;


public class CustomAwsSupplier implements AwsSecurityCredentialsSupplier {
    @Override
    public AwsSecurityCredentials getCredentials(ExternalAccountSupplierContext context)
            throws IOException {

         AwsCredentials dcc = WebIdentityTokenFileCredentialsProvider.builder().build().resolveCredentials();

         AwsSessionCredentials dc = (AwsSessionCredentials) dcc;

        return new AwsSecurityCredentials( dc.accessKeyId(), dc.secretAccessKey(),dc.sessionToken());
    }

    @Override
    public String getRegion(ExternalAccountSupplierContext context) throws IOException {
        return "us-east-2";
    }

}
import { expect as expectCDK, countResources } from '@aws-cdk/assert';
import * as cdk from '@aws-cdk/core';
import * as Permissions from '../lib/permissions-stack';

test('creates 2 managed policies', () => {
  const app = new cdk.App();
  const stack = new Permissions.AgcPermissionsStack(app, 'TestStack');

  expectCDK(stack).to(countResources("AWS::IAM::ManagedPolicy", 2))
});

test('admin policy is not empty', () => {
  const app = new cdk.App();
  const stack = new Permissions.AgcPermissionsStack(app, 'TestStack');

  let policy : string;
  policy = JSON.stringify(stack.adminPolicy.document.toJSON());
  expect(policy.length).toBeGreaterThan(0);
});

test('admin policy length within IAM limits', () => {
  // iam has a limit of 6144 characters per policy
  const app = new cdk.App();
  const stack = new Permissions.AgcPermissionsStack(app, 'TestStack');

  // CDK tokenizes Partition, Region, and AccountId with each synth so their
  // character contributions vary.
  // Substitute the longest possible values
  let policy : string;
  policy = JSON.stringify(stack.adminPolicy.document.toJSON())
    .replace(/\$\{Token\[AWS\.Partition\.\d+\]\}/g, "aws-us-gov")
    .replace(/\$\{Token\[AWS\.Region\.\d+\]\}/g, "ap-southwest-1")
    .replace(/\$\{Token\[AWS\.AccountId\.\d+\]\}/g, "444455556666");
  expect(policy.length).toBeLessThanOrEqual(6144);
});

test('user policy is not empty', () => {
  const app = new cdk.App();
  const stack = new Permissions.AgcPermissionsStack(app, 'TestStack');

  let policy : string;
  policy = JSON.stringify(stack.userPolicyPart1.document.toJSON());
  policy = JSON.stringify(stack.userPolicyPart2.document.toJSON());
  expect(policy.length).toBeGreaterThan(0);
});

test('user policy length within IAM limits', () => {
  // iam has a limit of 6144 characters per policy
  const app = new cdk.App();
  const stack = new Permissions.AgcPermissionsStack(app, 'TestStack');

  // CDK tokenizes Partition, Region, and AccountId with each synth so their
  // character contributions vary.
  // Substitute the longest possible values
  let policy : string;
  policy = JSON.stringify(stack.userPolicyPart1.document.toJSON())
    .replace(/\$\{Token\[AWS\.Partition\.\d+\]\}/g, "aws-us-gov")
    .replace(/\$\{Token\[AWS\.Region\.\d+\]\}/g, "ap-southwest-1")
    .replace(/\$\{Token\[AWS\.AccountId\.\d+\]\}/g, "444455556666");
  expect(policy.length).toBeLessThanOrEqual(6144);
  policy = JSON.stringify(stack.userPolicyPart2.document.toJSON())
      .replace(/\$\{Token\[AWS\.Partition\.\d+\]\}/g, "aws-us-gov")
      .replace(/\$\{Token\[AWS\.Region\.\d+\]\}/g, "ap-southwest-1")
      .replace(/\$\{Token\[AWS\.AccountId\.\d+\]\}/g, "444455556666");
  expect(policy.length).toBeLessThanOrEqual(6144);
});

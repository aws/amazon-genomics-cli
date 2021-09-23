#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from '@aws-cdk/core';
import { AgcPermissionsStack } from '../lib/permissions-stack';

const app = new cdk.App();
new AgcPermissionsStack(app, 'AgcPermissionsStack', {});

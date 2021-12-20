import { Node } from "constructs";
import { Maybe } from "../types";

export { CoreAppParameters } from "./core-app-parameters";
export { ContextAppParameters } from "./context-app-parameters";

export const getEnvString = (node: Node, key: string): string => {
  return getValue(node, key);
};

export const getEnvStringOrDefault = (node: Node, key: string, defaultValue?: string): Maybe<string> => {
  const value = node.tryGetContext(key);
  return value ? value : defaultValue;
};

export const getEnvBool = (node: Node, key: string): boolean => {
  return valueToBoolean(getValue(node, key));
};

export const getEnvBoolOrDefault = (node: Node, key: string, defaultValue?: boolean): Maybe<boolean> => {
  const value = node.tryGetContext(key);
  return value ? valueToBoolean(value) : defaultValue;
};

export const getEnvNumber = (node: Node, key: string): Maybe<number> => {
  const value = node.tryGetContext(key);
  return value ? Number(value) : undefined;
};

export const getEnvStringList = (node: Node, key: string): string[] => {
  return valueToList(getValue(node, key));
};

export const getEnvStringListOrDefault = (node: Node, key: string, defaultValue?: string[]): Maybe<string[]> => {
  const value = node.tryGetContext(key);
  return value ? valueToList(value) : defaultValue;
};

const getValue = (node: Node, key: string): string => {
  const value = node.tryGetContext(key);
  if (value === undefined || value === null || value === "") {
    throw Error(`App context cannot be null for key '${key}'`);
  }
  return value;
};

const valueToList = (value: string): string[] => {
  return value.split(",");
};

const valueToBoolean = (value: string): boolean => {
  return value === "true";
};

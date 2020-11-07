import axios from "axios";
import { notificationState } from "@/components/NotificationState";
import {NewPackage, Package} from "@/api/Models";

export interface Result {
  ID: number;
  Name: string;
  PackageBaseID: number;
  PackageBase: string;
  Version: string;
  Description: string;
  URL: string;
  NumVotes: number;
  Popularity: number;
  OutOfDate: null | number;
  Maintainer: string;
  FirstSubmitted: number;
  LastModified: number;
  URLPath: string;
}

export function ToPackage(result: Result | undefined): Package {

  if (typeof result === "undefined") {
    return NewPackage()
  }

  return {
    KeepLastN: 2,
    LastHash: [],
    Name: result.Name,
    RepoBranch: "master",
    RepoURL: `https://aur.archlinux.org/${result.Name}.git`,
    UpdateFrequency: 2 * 3600 * 1000 * 1000 * 1000
  }
}

export function NewResult(): Result {
  return {
    Description: "",
    FirstSubmitted: 0,
    ID: 0,
    LastModified: 0,
    Maintainer: "",
    Name: "",
    NumVotes: 0,
    OutOfDate: null,
    PackageBase: "",
    PackageBaseID: 0,
    Popularity: 0,
    URL: "",
    URLPath: "",
    Version: ""
  }
}

export interface Results {
  version: number;
  type: string;
  results: Result[];
  resultcount: number;
  error: undefined | string;
}

const client = axios.create({
  baseURL: " http://localhost:5001/search",
  headers: {
    "Content-type": "application/json"
  }
});

client.interceptors.response.use(undefined, error => {
  notificationState.message = error.message;
  notificationState.color = "#feb2b2";
  notificationState.enabled = true;

  return Promise.reject(error);
});

export async function search(term: string): Promise<Result[]> {
  if (term.length < 3) {
    return Promise.resolve([]);
  }

  return client.get("/" + term).then(resp => resp.data);
}

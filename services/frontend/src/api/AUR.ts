import axios from "axios";
import {notificationState} from "@/components/NotificationState";

export interface Result {
    ID: number,
    Name: string,
    PackageBaseID: number,
    PackageBase: string,
    Version: string,
    Description: string,
    URL: string,
    NumVotes: number,
    Popularity: number,
    OutOfDate: null | number,
    Maintainer: string,
    FirstSubmitted: number,
    LastModified: number,
    URLPath: string,
}

export interface Results {
    version: number,
    type: string,
    results: Result[],
    resultcount: number,
    error: undefined | string,
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
        return Promise.resolve([])
    }

    return client.get("/" + term).then(resp => resp.data)
}

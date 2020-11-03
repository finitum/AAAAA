import { GetPackages, Login} from "./API";
import { Package, User } from "@/api/Models";
import mockAxios from "jest-mock-axios";

afterEach(() => {
  // cleaning up the mess left behind the previous test
  mockAxios.reset();
});

describe("#API", () => {
  it("Should be able to handle a normal response", async () => {
    const promise = GetPackages();

    expect(mockAxios.get).toHaveBeenCalledWith("/package");

    const data: Package[] = [
      {
        Name: "string",
        RepoURL: "string",
        RepoBranch: "string",
        KeepLastN: 2,
        LastHash: "string",
        UpdateFrequency: 5
      }
    ];

    mockAxios.mockResponse({ data: data });

    const result = await promise;
    expect(result).toEqual(data);
  });
  it("Should be able to login", async () => {
    const u: User = {
      Username: "user",
      Password: "pass"
    };

    const promise = Login(u);

    expect(mockAxios.post).toHaveBeenCalledWith("/login", u);

    const data = { token: "token" };

    mockAxios.mockResponse({ data: data });

    const result = await promise;
    expect(result).toEqual(data.token);
  });
});

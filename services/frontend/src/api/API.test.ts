import {
  AddPackage,
  DeletePackage,
  GetAllUsers,
  GetPackages,
  Login,
  NewUser,
  UpdatePackage,
  UpdateUser
} from "./API";
import { Package, User } from "@/api/Models";
import mockAxios from "jest-mock-axios";

afterEach(() => {
  // cleaning up the mess left behind the previous test
  mockAxios.reset();
});

describe("#API", () => {
  it("Should be able to add a new user", async () => {
    const user: User = {
      Username: "test1",
      Password: "test2"
    };

    const promise = NewUser(user, "");

    expect(mockAxios.post).toHaveBeenCalledWith("/user", user);

    mockAxios.mockResponse({ data: {}, status: 200 });

    await promise;
  });

  it("Should be able to update a user", async () => {
    const user: User = {
      Username: "test1",
      Password: "test2"
    };

    const promise = UpdateUser(user, "");

    expect(mockAxios.put).toHaveBeenCalledWith("/user", user);

    mockAxios.mockResponse({ data: {}, status: 200 });

    await promise;
  });

  it("Should be able to get all users", async () => {
    const promise = GetAllUsers("");

    expect(mockAxios.get).toHaveBeenCalledWith("/users");

    mockAxios.mockResponse({ data: {}, status: 200 });

    await promise;
  });

  it("Should be able to update a package", async () => {
    const pkg: Package = {
      Name: "string",
      RepoURL: "string",
      RepoBranch: "string",
      KeepLastN: 2,
      UpdateFrequency: 5
    };

    const promise = UpdatePackage(pkg, "");

    expect(mockAxios.put).toHaveBeenCalledWith("/package/" + pkg.Name, pkg);

    mockAxios.mockResponse({ data: {}, status: 200 });

    await promise;
  });

  it("Should be able to delete a package", async () => {
    const pkgName = "something";

    const promise = DeletePackage(pkgName, "");

    expect(mockAxios.delete).toHaveBeenCalledWith("/package/" + pkgName);

    mockAxios.mockResponse({ data: {}, status: 200 });

    await promise;
  });

  it("Should be able to add a package", async () => {
    const pkg: Package = {
      Name: "string",
      RepoURL: "string",
      RepoBranch: "string",
      KeepLastN: 2,
      UpdateFrequency: 5
    };

    const promise = AddPackage(pkg, "");

    expect(mockAxios.post).toHaveBeenCalledWith("/package", pkg);

    mockAxios.mockResponse({ data: {}, status: 201 });

    await promise;
  });
  it("Should be able to handle a normal response", async () => {
    const promise = GetPackages();

    expect(mockAxios.get).toHaveBeenCalledWith("/package");

    const data: Package[] = [
      {
        Name: "string",
        RepoURL: "string",
        RepoBranch: "string",
        KeepLastN: 2,
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

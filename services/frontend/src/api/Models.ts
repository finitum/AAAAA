import humanizeDuration from "@/utils/humanizeDuration";

export interface Package {
  // Name is the name of the package (unique)
  Name: string;
  // RepoURL is the git repository where the PKGBUILD can be found
  RepoURL: string;
  // RepoBranch is which branch is used for updates
  RepoBranch: string;
  // KeepLastN determines how many old versions of packages are kept
  KeepLastN: number;
  // UpdateFrequency determines how often the package should be updated
  UpdateFrequency: number;
}

export function NewPackage(): Package {
  return {
    KeepLastN: 2,
    Name: "",
    RepoBranch: "",
    RepoURL: "",
    UpdateFrequency: 0
  };
}

export function frequencyToDuration(freqns: number): string {
  return humanizeDuration(freqns / 1000 / 1000);
}

export interface User {
  Username: string;
  Password: string;
}

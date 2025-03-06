import axios from "axios";

const serverAPIURL = "/api/ui";

function publishEvent401() {
  const eventAwesome = new CustomEvent("noauth", {
    bubbles: true,
  });
  document.dispatchEvent(eventAwesome);
}

export async function fetchBookmarks(limit: Number, nextID: Number, fetchType: string): Promise<any> {
  const apiURL = `${serverAPIURL}/url/all-bookmarks?limit=${limit}&next_id=${nextID}&fetch_type=${fetchType}`;
  return new Promise((resolve, reject) => {
    axios
      .get(apiURL)
      .then((response) => {
        resolve(response.data);
      })
      .catch((error) => {
        //check if status code is 401
        if (error.response && error.response.status === 401) publishEvent401();
        reject(error);
      });
  });
}

export async function callSearch(needle: String): Promise<any> {
  const apiURL = `${serverAPIURL}/url/search-bookmarks`;
  const body = JSON.stringify({ needle });
  const headers = {
    "Content-Type": "application/json",
  };
  try {
    const response = await fetch(apiURL, { method: "POST", headers, body });
    const data = await response.json();
    return data;
  } catch (error) {
    console.error("Error searching bookmarks:", error);
    throw error;
  }
}
export async function addNewBookmarkURL(url: String): Promise<any> {
  const apiURL = `${serverAPIURL}/url/new-bookmark`;
  const body = JSON.stringify({ url });
  const headers = {
    "Content-Type": "application/json",
  };

  //use axios
  return new Promise((resolve, reject) => {
    axios
      .post(apiURL, body, { headers })
      .then((response) => {
        resolve(response.data);
      })
      .catch((error) => {
        reject(error);
      });
  });
}
export async function fetchGithubStarredRepos(username: String): Promise<any> {
  const apiURL = `${serverAPIURL}/url/import-github`;
  const body = JSON.stringify({ username });
  const headers = {
    "Content-Type": "application/json",
  };

  //use axios
  return new Promise((resolve, reject) => {
    axios
      .post(apiURL, body, { headers })
      .then((response) => {
        resolve(response.data);
      })
      .catch((error) => {
        reject(error);
      });
  });
}

export async function bulkAddUrls(urls: String[], direction: String): Promise<any> {
  const apiURL = `${serverAPIURL}/url/import-bulk`;
  const body = JSON.stringify({ urls, direction });
  const headers = {
    "Content-Type": "application/json",
  };

  //use axios
  return new Promise((resolve, reject) => {
    axios
      .post(apiURL, body, { headers })
      .then((response) => {
        resolve(response.data);
      })
      .catch((error) => {
        reject(error);
      });
  });
}

export async function countBookmarks(): Promise<any> {
  const apiURL = `${serverAPIURL}/url/get-bookmark-count`;

  return new Promise((resolve, reject) => {
    axios
      .get(apiURL)
      .then((response) => {
        resolve(response.data);
      })
      .catch((error) => {
        reject(error);
      });
  });
}

export async function uploadFirefoxFile(file: File): Promise<any> {
  const apiURL = `${serverAPIURL}/url/import-browsers`;
  const formData = new FormData();
  formData.append("file", file);

  //use axios
  return new Promise((resolve, reject) => {
    axios
      .post(apiURL, formData, {
        headers: {
          "Content-Type": "multipart/form-data",
        },
      })
      .then((response) => {
        resolve(response.data);
      })
      .catch((error) => {
        reject(error);
      });
  });
}

export async function getUser(): Promise<any> {
  const apiURL = `${serverAPIURL}/user`;

  return new Promise((resolve, reject) => {
    axios
      .get(apiURL)
      .then((response) => {
        resolve(response.data);
      })
      .catch((error) => {
        reject(error);
      });
  });
}

export async function getBookmarkById(id: String): Promise<any> {
  const apiURL = `${serverAPIURL}/url/get-bookmark/${id}`;
  return new Promise((resolve, reject) => {
    axios
      .get(apiURL)
      .then((response) => {
        resolve(response.data);
      })
      .catch((error) => {
        reject(error);
      });
  });
}
export async function deleteBookmarkById(id: String): Promise<any> {
  const apiURL = `${serverAPIURL}/url/delete-bookmark/${id}`;
  return new Promise((resolve, reject) => {
    axios
      .delete(apiURL)
      .then((response) => {
        resolve(response.data);
      })
      .catch((error) => {
        reject(error);
      });
  });
}
export async function getJobQueueStatus(): Promise<any> {
  const apiURL = `${serverAPIURL}/url/bookmarks-queue`;
  return new Promise((resolve, reject) => {
    axios
      .get(apiURL)
      .then((response) => {
        resolve(response.data);
      })
      .catch((error) => {
        reject(error);
      });
  });
}

//implement patch
export async function updateBookmark(id: String, data: any): Promise<any> {
  const apiURL = `${serverAPIURL}/url/bookmark-update/${id}`;
  const body = JSON.stringify(data);
  const headers = {
    "Content-Type": "application/json",
  };

  //use axios
  return new Promise((resolve, reject) => {
    axios
      .patch(apiURL, body, { headers })
      .then((response) => {
        resolve(response.data);
      })
      .catch((error) => {
        reject(error);
      });
  });
}
export async function refetchBookmark(id: String): Promise<any> {
  const apiURL = `${serverAPIURL}/url/index-bookmark/${id}`;
  const headers = {
    "Content-Type": "application/json",
  };

  //use axios
  return new Promise((resolve, reject) => {
    axios
      .patch(apiURL, null, { headers })
      .then((response) => {
        resolve(response.data);
      })
      .catch((error) => {
        reject(error);
      });
  });
}

//do a post call to /api/bookmarks/delete/bulk
export async function deleteBulkBookmarks(ids: String[]): Promise<any> {
  const apiURL = `${serverAPIURL}/url/delete-bulk`;
  const body = JSON.stringify({ organization_relation_ids: ids });
  const headers = {
    "Content-Type": "application/json",
  };
  return new Promise((resolve, reject) => {
    axios
      .post(apiURL, body, { headers })
      .then((response) => {
        resolve(response.data);
      })
      .catch((error) => {
        reject(error);
      });
  });
}

export async function getExtensionSecrets(): Promise<any> {
  const apiURL = `${serverAPIURL}/secrets/extensions/active`;
  return new Promise((resolve, reject) => {
    axios
      .get(apiURL)
      .then((response) => {
        resolve(response.data);
      })
      .catch((error) => {
        reject(error);
      });
  });
}

export async function createNewSecret(secretName: String): Promise<any> {
  const apiURL = `${serverAPIURL}/secrets/extensions`;
  const body = JSON.stringify({ secret_name: secretName });
  const headers = {
    "Content-Type": "application/json",
  };
  return new Promise((resolve, reject) => {
    axios
      .post(apiURL, body, { headers })
      .then((response) => {
        resolve(response.data);
      })
      .catch((error) => {
        reject(error);
      });
  });
}
export async function deactivateSecret(secretID: String): Promise<any> {
  const apiURL = `${serverAPIURL}/secrets/deactivate`;
  const body = JSON.stringify({ secret_id: secretID });
  const headers = {
    "Content-Type": "application/json",
  };
  return new Promise((resolve, reject) => {
    axios
      .post(apiURL, body, { headers })
      .then((response) => {
        resolve(response.data);
      })
      .catch((error) => {
        reject(error);
      });
  });
}

export async function getMySubscription(): Promise<any> {
  const apiURL = `${serverAPIURL}/billing/details`;
  return new Promise((resolve, reject) => {
    axios
      .get(apiURL)
      .then((response) => {
        resolve(response.data);
      })
      .catch((error) => {
        reject(error);
      });
  });
}

export async function login(email: string, password: string): Promise<any> {
  const apiURL = `${serverAPIURL}/user/login`;
  const body = JSON.stringify({ email, password });
  const headers = {
    "Content-Type": "application/json",
  };
  return new Promise((resolve, reject) => {
    axios
      .post(apiURL, body, { headers })
      .then((response) => {
        resolve(response.data);
      })
      .catch((error) => {
        reject(error);
      });
  });
}

export async function signup(email: string, password: string, name: string): Promise<any> {
  const apiURL = `${serverAPIURL}/user/sign-up`;
  const body = JSON.stringify({ email, password, name });
  const headers = {
    "Content-Type": "application/json",
  };
  return new Promise((resolve, reject) => {
    axios
      .post(apiURL, body, { headers })
      .then((response) => {
        resolve(response.data);
      })
      .catch((error) => {
        reject(error);
      });
  });
}

export async function getSchedules(id: String): Promise<any> {
  const apiURL = `${serverAPIURL}/url/view-schedules`;
  return new Promise((resolve, reject) => {
    axios
      .get(apiURL)
      .then((response) => {
        resolve(response.data);
      })
      .catch((error) => {
        reject(error);
      });
  });
}

export async function updateSchedule(schedule_id: String, status: String, interval: Number): Promise<any> {
  const apiURL = `${serverAPIURL}/url/update-schedules`;
  const headers = {
    "Content-Type": "application/json",
  };
  const body = JSON.stringify({ schedule_id, status, interval });
  //use axios
  return new Promise((resolve, reject) => {
    axios
      .patch(apiURL, body, { headers })
      .then((response) => {
        resolve(response.data);
      })
      .catch((error) => {
        reject(error);
      });
  });
}

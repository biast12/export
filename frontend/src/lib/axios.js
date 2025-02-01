import {PUBLIC_BACKEND_URI} from "$env/static/public";
import axios from "axios";
import {goto} from "$app/navigation";

export const client = axios.create({
  baseURL: PUBLIC_BACKEND_URI,
  headers: {
    'Authorization': `Bearer ${window.localStorage.getItem('token')}`,
  },
  validateStatus: false,
});

client.interceptors.response.use((res) => {
    if (res.status === 401) {
        window.localStorage.clear();
        goto("/auth/sign-in");
    }

    return res;
})


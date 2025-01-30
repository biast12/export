import {goto} from "$app/navigation";

export function load() {
    if (!window.localStorage.getItem("token")) {
        goto("/auth/sign-in");
    }
}
import {Auth0AuthProvider} from "ra-auth-auth0";
import {Auth0Client} from "@auth0/auth0-spa-js";

export const auth0 = new Auth0Client({
    domain: "dev-ip4wfrm3ukv6cdnc.us.auth0.com",
    clientId: "t5RWIxHNCc9STAWcGB5i2wfiKJQNDhCt",
    cacheLocation: 'localstorage',
    authorizationParams: {
        redirect_uri: window.location.origin
    },
});

export const authProvider = Auth0AuthProvider(auth0, {
    loginRedirectUri: window.location.origin + "/#/auth-callback",
    logoutRedirectUri: window.location.origin
})

export default authProvider

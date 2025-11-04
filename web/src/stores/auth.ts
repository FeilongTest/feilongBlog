import { ref } from "vue";
import { defineStore } from "pinia";
import ApiService from "@/core/services/ApiService";
import JwtService from "@/core/services/JwtService";
import service from "@/utils/request"

export interface UserData {
  token:string;
  expiresAt:number;
  user:User;
}

export interface User {
  ID:number;
  userName:string;
  pic:string;
  email:string;
  ctime:number;
  lasttime:number;
  ip:string;
  status:number;
  truename:string;
  admin:number;
}


export const useAuthStore = defineStore("auth", () => {
  const errors = ref({});
  const user = ref<UserData>({} as UserData);
  const isAuthenticated = ref(!!JwtService.getToken());

  function setAuth(authUser: UserData) {
    isAuthenticated.value = true;
    user.value = authUser;
    errors.value = {};
    JwtService.saveToken(user.value.token);
  }

  function setError(error: any) {
    errors.value = { ...error };
  }

  function purgeAuth() {
    isAuthenticated.value = false;
    user.value = {} as UserData;
    errors.value = [];
    JwtService.destroyToken();
  }

  function login(credentials: User) {
    return service.post("/base/login",credentials)
    .then(({data})=>{
      if(data == undefined || data.token == ""){
        setError(data.msg)
        return
      }
      setAuth(data)
    })
    .catch(({ response}) => {
      setError(response.data.errors);
    })
  }

  function logout() {
    purgeAuth();
  }

  function register(credentials: User) {
    return ApiService.post("register", credentials)
      .then(({ data }) => {
        setAuth(data);
      })
      .catch(({ response }) => {
        setError(response.data.errors);
      });
  }

  function forgotPassword(email: string) {
    return ApiService.post("forgot_password", email)
      .then(() => {
        setError({});
      })
      .catch(({ response }) => {
        setError(response.data.errors);
      });
  }

  function verifyAuth() {
    // if (JwtService.getToken()) {
    //   service.post("verify_token", { token: JwtService.getToken() })
    //     .then(({ data }) => {
    //       setAuth(data);
    //     })
    //     .catch(({ response }) => {
    //       setError(response.data.errors);
    //       purgeAuth();
    //     });
    // } else {
    //   purgeAuth();
    // }
  }

  return {
    errors,
    user,
    isAuthenticated,
    login,
    logout,
    register,
    forgotPassword,
    verifyAuth,
  };
});

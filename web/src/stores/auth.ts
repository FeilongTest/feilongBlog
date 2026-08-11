import { ref } from "vue";
import { defineStore } from "pinia";
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
  bio:string;
  admin:number;
}

export interface PublicProfile {
  trueName:string;
  email:string;
  pic:string;
  bio:string;
}

let publicProfileRequest: Promise<PublicProfile> | null = null;


export const useAuthStore = defineStore("auth", () => {
  const errors = ref({});
  const user = ref<UserData>({} as UserData);
  const isAuthenticated = ref(JwtService.isTokenValid());
  const publicProfile = ref<PublicProfile>({} as PublicProfile);

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
	.catch((error: any) => {
	  setError({ login: error?.response?.data?.msg || error?.message || "登录失败" });
    })
  }

  function logout() {
    purgeAuth();
  }

  function updateProfile(profile: User) {
    if (user.value?.user) user.value.user = profile;
    else user.value = { token: JwtService.getToken() || "", expiresAt: 0, user: profile };
  }

  async function refreshProfile() {
    if (!isAuthenticated.value) return undefined;
    const { data } = await service.get("/admin/user/profile");
    updateProfile(data);
    return data as User;
  }

  async function refreshPublicProfile() {
    if (publicProfile.value?.email || publicProfile.value?.bio || publicProfile.value?.pic) {
      return publicProfile.value;
    }
    if (!publicProfileRequest) {
      publicProfileRequest = service.get("/base/getContact")
        .then(({ data }) => {
          publicProfile.value = data || {};
          return publicProfile.value;
        })
        .finally(() => {
          publicProfileRequest = null;
        });
    }
    return publicProfileRequest;
  }

  function verifyAuth() {
	  isAuthenticated.value = JwtService.isTokenValid();
	  if (!isAuthenticated.value) JwtService.destroyToken();
  }

  return {
    errors,
    user,
    publicProfile,
    isAuthenticated,
    login,
    logout,
    updateProfile,
    refreshProfile,
    refreshPublicProfile,
    verifyAuth,
  };
});

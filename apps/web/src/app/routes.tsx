import { Navigate, RouteObject } from "react-router-dom";
import { FeedPage } from "../features/feed/pages/FeedPage";
import { LoginPage } from "../features/auth/pages/LoginPage";
import { SignupPage } from "../features/auth/pages/SignupPage";
import { MyProfilePage } from "../features/profile/pages/MyProfilePage";
import { PublicProfilePage } from "../features/profile/pages/PublicProfilePage";
export const appRoutes:RouteObject[]=[{path:"/",element:<Navigate to="/feed" replace/>},{path:"/login",element:<LoginPage/>},{path:"/signup",element:<SignupPage/>},{path:"/feed",element:<FeedPage/>},{path:"/profile/me",element:<MyProfilePage/>},{path:"/profile/:userId",element:<PublicProfilePage/>}];

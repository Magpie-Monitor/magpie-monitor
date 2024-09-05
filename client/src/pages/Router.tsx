import {
  Route,
  createBrowserRouter,
  createRoutesFromElements,
  defer,
} from "react-router-dom";
import Login from "./Login/Login";
import { ProtectedLayout } from "../providers/AuthProvider/ProtectedLayout";
import { AuthLayout } from "../providers/AuthProvider/AuthLayout";
import { getTokenInfo } from "../api/authApi";

const getUserData = () => {
  return getTokenInfo();
};

const router = createBrowserRouter(
  createRoutesFromElements(
    <Route
      element={<AuthLayout />}
      loader={() => {
        return defer({
          userData: getUserData(),
        });
      }}
      errorElement={<Login />}
    >
      <Route path="/" element={<Login />} />
      <Route path="/" element={<ProtectedLayout />}></Route>
    </Route>,
  ),
);

export default router;

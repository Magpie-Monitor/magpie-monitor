import { Route, createBrowserRouter, createRoutesFromElements, defer } from 'react-router-dom';
import Login from './Login/Login';
// import { ProtectedLayout } from 'providers/AuthProvider/ProtectedLayout';
import { AuthLayout } from 'providers/AuthProvider/AuthLayout';
import Home from './Home/Home';
import { getAuthInfo } from 'providers/AuthProvider/AuthProvider';

const router = createBrowserRouter(
  createRoutesFromElements(
    <Route
      element={<AuthLayout />}
      loader={() => {
        return defer({
          authData: getAuthInfo(),
        });
      }}
      errorElement={<Home />}
    >
      <Route path="/login" element={<Login />} />
      <Route path="/" element={<Home />} />
    </Route>,
  ),
);

// <Route path="/" element={<ProtectedLayout />}>
// </Route>

export default router;

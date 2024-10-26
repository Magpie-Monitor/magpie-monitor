import {
  Route,
  createBrowserRouter,
  createRoutesFromElements,
  defer,
} from 'react-router-dom';
import Login from './Login/Login';
import { ProtectedLayout } from 'providers/AuthProvider/ProtectedLayout';
import { AuthLayout } from 'providers/AuthProvider/AuthLayout';
import Home from './Home/Home';
import { getAuthInfo } from 'providers/AuthProvider/AuthProvider';
import NotFoundError from './NotFoundError/NotFoundError';
import Reports from './Reports/Reports.tsx';
import Settings from './Settings/Settings.tsx';
import ReportDetails from './ReportDetails/ReportDetails.tsx';

const router = createBrowserRouter(
  createRoutesFromElements(
    <Route
      element={<AuthLayout />}
      loader={() => {
        return defer({
          authData: getAuthInfo(),
        });
      }}
      errorElement={<NotFoundError />}
    >
      <Route path="/login" element={<Login />} />
      <Route path="/" element={<ProtectedLayout />}>
        <Route path="/" element={<Home />} />
        <Route path="/dashboard" element={<Home />} />
        <Route path="/reports" >
          <Route path="" element={<Reports />} />
          <Route path=":id" element={<ReportDetails />} />
        </Route>
        <Route path="/settings" element={<Settings />} />
      </Route>
    </Route>,
  ),
);

export default router;

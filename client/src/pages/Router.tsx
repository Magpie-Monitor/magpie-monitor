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
import ReportDetails from './ReportDetails/ReportDetails.tsx';
import Notification from './Notification/Notification.tsx';
import Clusters from './Clusters/Clusters.tsx';
import NodeIncidentPage from './Incident/NodeIncident.tsx';
import ApplicationIncidentPage from './Incident/ApplicationIncident.tsx';
import CreateReport from './Report/CreateReport.tsx';

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
        <Route path="/clusters">
          <Route path="" element={<Clusters />} />
          <Route path=":id/report" element={<CreateReport />} />
        </Route>
        <Route path="/node-incidents/:id" element={<NodeIncidentPage />} />
        <Route
          path="/application-incidents/:id"
          element={<ApplicationIncidentPage />}
        />
        <Route path="/reports">
          <Route path="" element={<Reports />} />
          <Route path=":id" element={<ReportDetails />} />
        </Route>
        <Route path="/notifications" element={<Notification />} />
      </Route>
    </Route>,
  ),
);

export default router;

import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.scss';
import Router from 'pages/Router.tsx';
import { RouterProvider } from 'react-router-dom';
import ChartJS from 'chart.js/auto';
import { LinearScale, TimeScale } from 'chart.js';
import { ToastProvider } from 'providers/ToastProvider/ToastProvider';

ChartJS.register(TimeScale, LinearScale);

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <ToastProvider>
      <RouterProvider router={Router} />
    </ToastProvider>
  </React.StrictMode>,
);

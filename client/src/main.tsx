import React from "react";
import ReactDOM from "react-dom/client";
import "./index.css";
import Router from "./pages/Router.tsx";
import { RouterProvider } from "react-router-dom";

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <RouterProvider router={Router} />
  </React.StrictMode>,
);

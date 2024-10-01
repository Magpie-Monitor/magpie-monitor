import { useEffect, useState } from 'react';
import { useAuth } from 'providers/AuthProvider/AuthProvider';
import { login } from 'api/authApi';
import { useNavigate } from 'react-router-dom';

const getCodeFromParams = () => {
  const queryParams = new URLSearchParams(window.location.search);
  return queryParams.get('code');
};

const useLogin = () => {
  const [code] = useState(getCodeFromParams());
  const { authenticationInfo, setAuthenticationInfo, isTokenValid } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    if (code) {
      login(code)
        .then((resolvedAuthenticationInfo) => {
          setAuthenticationInfo(resolvedAuthenticationInfo);
          navigate('/config');
        })
        .catch((error) => {
          console.error('Failed to login', error); // eslint-disable-line no-console
        });
    }

    if (isTokenValid()) {
      navigate('/config');
    }
  }, [code, setAuthenticationInfo, navigate, isTokenValid]);

  return authenticationInfo;
};

export default useLogin;

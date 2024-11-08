import { useEffect } from 'react';
import { useAuth } from 'providers/AuthProvider/AuthProvider';
import { useNavigate } from 'react-router-dom';

const useLogin = () => {
  const { authenticationInfo, setAuthenticationInfo, isTokenValid } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    isTokenValid().then((isValid) => {
      if (isValid) {
        navigate('/dashboard');
      }
    });
  }, [setAuthenticationInfo, navigate, isTokenValid]);

  return authenticationInfo;
};

export default useLogin;

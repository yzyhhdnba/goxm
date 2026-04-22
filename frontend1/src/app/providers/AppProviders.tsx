import { PropsWithChildren, useEffect } from 'react';
import { Provider, useDispatch } from 'react-redux';
import { BrowserRouter } from 'react-router-dom';
import { registerUnauthorizedHandler } from '@/api/interceptors';
import { bootstrapAuth, forceLogout } from '@/store/slices/authSlice';
import { AppDispatch, store } from '@/store';

function Bootstrapper({ children }: PropsWithChildren) {
  const dispatch = useDispatch<AppDispatch>();

  useEffect(() => {
    dispatch(bootstrapAuth());
  }, [dispatch]);

  useEffect(() => {
    registerUnauthorizedHandler(() => {
      dispatch(forceLogout('登录状态已失效，请重新登录'));
    });

    return () => {
      registerUnauthorizedHandler(null);
    };
  }, [dispatch]);

  return children;
}

export function AppProviders({ children }: PropsWithChildren) {
  return (
    <Provider store={store}>
      <BrowserRouter>
        <Bootstrapper>{children}</Bootstrapper>
      </BrowserRouter>
    </Provider>
  );
}

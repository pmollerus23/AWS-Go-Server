import { useState, useEffect, useCallback, useRef } from 'react';
import type { QueryState, ApiError, AsyncFunction } from '../types';

interface UseQueryOptions<T> {
  enabled?: boolean;
  refetchOnMount?: boolean;
  refetchInterval?: number;
  onSuccess?: (data: T) => void;
  onError?: (error: ApiError) => void;
  retry?: number;
  retryDelay?: number;
}

interface UseQueryResult<T> extends QueryState<T> {
  refetch: () => Promise<void>;
  reset: () => void;
}

export const useQuery = <T>(
  queryFn: AsyncFunction<T>,
  options: UseQueryOptions<T> = {}
): UseQueryResult<T> => {
  const {
    enabled = true,
    refetchOnMount = true,
    refetchInterval,
    onSuccess,
    onError,
    retry = 0,
    retryDelay = 1000,
  } = options;

  const [state, setState] = useState<QueryState<T>>({
    data: null,
    error: null,
    status: 'idle',
    isLoading: false,
    isSuccess: false,
    isError: false,
  });

  const retryCountRef = useRef(0);
  const isMountedRef = useRef(true);
  const onSuccessRef = useRef(onSuccess);
  const onErrorRef = useRef(onError);

  // Update refs when callbacks change
  useEffect(() => {
    onSuccessRef.current = onSuccess;
  }, [onSuccess]);

  useEffect(() => {
    onErrorRef.current = onError;
  }, [onError]);

  const executeQuery = useCallback(async (): Promise<void> => {
    if (!enabled) return;

    setState(prev => ({
      ...prev,
      status: 'loading',
      isLoading: true,
      isError: false,
    }));

    try {
      const data = await queryFn();

      if (!isMountedRef.current) return;

      setState({
        data,
        error: null,
        status: 'success',
        isLoading: false,
        isSuccess: true,
        isError: false,
      });

      retryCountRef.current = 0;
      onSuccessRef.current?.(data);
    } catch (err) {
      if (!isMountedRef.current) return;

      const error: ApiError = {
        message: err instanceof Error ? err.message : 'An error occurred',
        status: (err as any)?.status,
        code: (err as any)?.code,
      };

      // Retry logic
      if (retryCountRef.current < retry) {
        retryCountRef.current++;
        setTimeout(() => {
          executeQuery();
        }, retryDelay);
        return;
      }

      setState({
        data: null,
        error,
        status: 'error',
        isLoading: false,
        isSuccess: false,
        isError: true,
      });

      retryCountRef.current = 0;
      onErrorRef.current?.(error);
    }
  }, [queryFn, enabled, retry, retryDelay]);

  const refetch = useCallback(async (): Promise<void> => {
    retryCountRef.current = 0;
    await executeQuery();
  }, [executeQuery]);

  const reset = useCallback((): void => {
    setState({
      data: null,
      error: null,
      status: 'idle',
      isLoading: false,
      isSuccess: false,
      isError: false,
    });
    retryCountRef.current = 0;
  }, []);

  // Initial fetch and refetch on mount
  useEffect(() => {
    if (refetchOnMount && enabled) {
      executeQuery();
    }
  }, [executeQuery, refetchOnMount, enabled]);

  // Refetch interval
  useEffect(() => {
    if (!refetchInterval || !enabled) return;

    const intervalId = setInterval(() => {
      executeQuery();
    }, refetchInterval);

    return () => clearInterval(intervalId);
  }, [refetchInterval, enabled, executeQuery]);

  // Track mount/unmount status
  useEffect(() => {
    isMountedRef.current = true;
    return () => {
      isMountedRef.current = false;
    };
  }, []);

  return {
    ...state,
    refetch,
    reset,
  };
};

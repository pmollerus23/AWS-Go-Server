import { useState, useCallback, useRef, useEffect } from 'react';
import type { ApiError } from '../types';

interface MutationState<T> {
  data: T | null;
  error: ApiError | null;
  isLoading: boolean;
  isSuccess: boolean;
  isError: boolean;
}

interface UseMutationOptions<T, TVariables> {
  onSuccess?: (data: T, variables: TVariables) => void;
  onError?: (error: ApiError, variables: TVariables) => void;
  onSettled?: (data: T | null, error: ApiError | null, variables: TVariables) => void;
}

interface UseMutationResult<T, TVariables> extends MutationState<T> {
  mutate: (variables: TVariables) => Promise<void>;
  mutateAsync: (variables: TVariables) => Promise<T>;
  reset: () => void;
}

export const useMutation = <T, TVariables = void>(
  mutationFn: (variables: TVariables) => Promise<T>,
  options: UseMutationOptions<T, TVariables> = {}
): UseMutationResult<T, TVariables> => {
  const { onSuccess, onError, onSettled } = options;

  const [state, setState] = useState<MutationState<T>>({
    data: null,
    error: null,
    isLoading: false,
    isSuccess: false,
    isError: false,
  });

  // Use refs for callbacks to avoid recreating mutateAsync on every render
  const onSuccessRef = useRef(onSuccess);
  const onErrorRef = useRef(onError);
  const onSettledRef = useRef(onSettled);

  // Update refs when callbacks change
  useEffect(() => {
    onSuccessRef.current = onSuccess;
  }, [onSuccess]);

  useEffect(() => {
    onErrorRef.current = onError;
  }, [onError]);

  useEffect(() => {
    onSettledRef.current = onSettled;
  }, [onSettled]);

  const mutateAsync = useCallback(
    async (variables: TVariables): Promise<T> => {
      setState({
        data: null,
        error: null,
        isLoading: true,
        isSuccess: false,
        isError: false,
      });

      try {
        const data = await mutationFn(variables);

        setState({
          data,
          error: null,
          isLoading: false,
          isSuccess: true,
          isError: false,
        });

        onSuccessRef.current?.(data, variables);
        onSettledRef.current?.(data, null, variables);

        return data;
      } catch (err) {
        const error: ApiError = {
          message: err instanceof Error ? err.message : 'An error occurred',
          status: (err as any)?.status,
          code: (err as any)?.code,
        };

        setState({
          data: null,
          error,
          isLoading: false,
          isSuccess: false,
          isError: true,
        });

        onErrorRef.current?.(error, variables);
        onSettledRef.current?.(null, error, variables);

        throw error;
      }
    },
    [mutationFn]
  );

  const mutate = useCallback(
    async (variables: TVariables): Promise<void> => {
      try {
        await mutateAsync(variables);
      } catch {
        // Error already handled in mutateAsync
      }
    },
    [mutateAsync]
  );

  const reset = useCallback((): void => {
    setState({
      data: null,
      error: null,
      isLoading: false,
      isSuccess: false,
      isError: false,
    });
  }, []);

  return {
    ...state,
    mutate,
    mutateAsync,
    reset,
  };
};

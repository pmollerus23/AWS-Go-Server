import type { ReactNode } from 'react';

export type PropsWithChildren<P = unknown> = P & {
  children?: ReactNode;
};

export type Nullable<T> = T | null;

export type Optional<T> = T | undefined;

export type AsyncFunction<T = void, Args extends unknown[] = []> = (
  ...args: Args
) => Promise<T>;

export type VoidFunction = () => void;

export interface BaseComponentProps {
  className?: string;
  id?: string;
  testId?: string;
}

export type DeepPartial<T> = {
  [P in keyof T]?: T[P] extends object ? DeepPartial<T[P]> : T[P];
};

export type RequireAtLeastOne<T, Keys extends keyof T = keyof T> = Pick<
  T,
  Exclude<keyof T, Keys>
> &
  {
    [K in Keys]-?: Required<Pick<T, K>> & Partial<Pick<T, Exclude<Keys, K>>>;
  }[Keys];

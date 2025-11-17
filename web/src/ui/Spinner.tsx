import type { BaseComponentProps } from '../types';

interface SpinnerProps extends BaseComponentProps {
  size?: 'small' | 'medium' | 'large';
  color?: 'primary' | 'secondary' | 'white';
  centered?: boolean;
}

export const Spinner: React.FC<SpinnerProps> = ({
  size = 'medium',
  color = 'primary',
  centered = false,
  className = '',
  id,
  testId,
}) => {
  const classes = [
    'spinner',
    `spinner-${size}`,
    `spinner-${color}`,
    className,
  ]
    .filter(Boolean)
    .join(' ');

  const spinner = (
    <div className={classes} id={id} data-testid={testId}>
      <div className="spinner-circle" />
    </div>
  );

  if (centered) {
    return <div className="spinner-container-centered">{spinner}</div>;
  }

  return spinner;
};

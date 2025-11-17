import type { PropsWithChildren, BaseComponentProps } from '../types';

interface CardProps extends BaseComponentProps, PropsWithChildren {
  title?: string;
  footer?: React.ReactNode;
  padding?: 'none' | 'small' | 'medium' | 'large';
  elevation?: 'none' | 'low' | 'medium' | 'high';
}

export const Card: React.FC<CardProps> = ({
  children,
  title,
  footer,
  padding = 'medium',
  elevation = 'medium',
  className = '',
  id,
  testId,
}) => {
  const classes = [
    'card',
    `card-padding-${padding}`,
    `card-elevation-${elevation}`,
    className,
  ]
    .filter(Boolean)
    .join(' ');

  return (
    <div className={classes} id={id} data-testid={testId}>
      {title && (
        <div className="card-header">
          <h3 className="card-title">{title}</h3>
        </div>
      )}

      <div className="card-content">{children}</div>

      {footer && <div className="card-footer">{footer}</div>}
    </div>
  );
};

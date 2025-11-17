import { Card } from '../ui';
import type { User, BaseComponentProps } from '../types';

interface UserCardProps extends BaseComponentProps {
  user: User;
  onEdit?: () => void;
}

export const UserCard: React.FC<UserCardProps> = ({
  user,
  onEdit,
  className,
  id,
  testId,
}) => {
  const footer = onEdit ? (
    <button onClick={onEdit} className="btn btn-secondary">
      Edit Profile
    </button>
  ) : undefined;

  return (
    <Card
      title="User Profile"
      footer={footer}
      className={className}
      id={id}
      testId={testId}
    >
      <div className="user-card-content">
        <div className="user-details">
          {user.name && (
            <p>
              <strong>Name:</strong> {user.name}
            </p>
          )}
          {(user.firstName || user.lastName) && (
            <p>
              <strong>Full Name:</strong> {`${user.firstName || ''} ${user.lastName || ''}`.trim()}
            </p>
          )}
          <p>
            <strong>Username:</strong> {user.username}
          </p>
          <p>
            <strong>Email:</strong> {user.email}
          </p>
          {user.groups && user.groups.length > 0 && (
            <p>
              <strong>Groups:</strong> {user.groups.join(', ')}
            </p>
          )}
        </div>
      </div>
    </Card>
  );
};

import { useState } from 'react';
import { useAuth } from '../contexts';
import { useMutation } from '../hooks';
import type { User } from '../types';

export const ProfilePage: React.FC = () => {
  const { user, updateUser } = useAuth();
  const [isEditing, setIsEditing] = useState(false);
  const [formData, setFormData] = useState({
    firstName: user?.firstName || '',
    lastName: user?.lastName || '',
    email: user?.email || '',
  });

  const updateProfileMutation = useMutation<User, Partial<User>>(
    async (data) => {
      // TODO: Replace with actual API call
      // return authApi.updateProfile(data);
      await new Promise(resolve => setTimeout(resolve, 1000));
      return { ...user!, ...data };
    },
    {
      onSuccess: (data) => {
        updateUser(data);
        setIsEditing(false);
      },
      onError: (error) => {
        console.error('Failed to update profile:', error);
      },
    }
  );

  const handleSubmit = async (e: React.FormEvent): Promise<void> => {
    e.preventDefault();
    await updateProfileMutation.mutate(formData);
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>): void => {
    setFormData(prev => ({
      ...prev,
      [e.target.name]: e.target.value,
    }));
  };

  if (!user) {
    return (
      <div className="profile-page">
        <p>Please log in to view your profile.</p>
      </div>
    );
  }

  return (
    <div className="profile-page">
      <h1>Profile</h1>

      {!isEditing ? (
        <div className="profile-view">
          <div className="profile-field">
            <label>First Name:</label>
            <span>{user.firstName || 'Not set'}</span>
          </div>

          <div className="profile-field">
            <label>Last Name:</label>
            <span>{user.lastName || 'Not set'}</span>
          </div>

          <div className="profile-field">
            <label>Email:</label>
            <span>{user.email}</span>
          </div>

          <div className="profile-field">
            <label>Username:</label>
            <span>{user.username}</span>
          </div>

          <button onClick={() => setIsEditing(true)} className="edit-button">
            Edit Profile
          </button>
        </div>
      ) : (
        <form onSubmit={handleSubmit} className="profile-form">
          <div className="form-field">
            <label htmlFor="firstName">First Name:</label>
            <input
              type="text"
              id="firstName"
              name="firstName"
              value={formData.firstName}
              onChange={handleChange}
            />
          </div>

          <div className="form-field">
            <label htmlFor="lastName">Last Name:</label>
            <input
              type="text"
              id="lastName"
              name="lastName"
              value={formData.lastName}
              onChange={handleChange}
            />
          </div>

          <div className="form-field">
            <label htmlFor="email">Email:</label>
            <input
              type="email"
              id="email"
              name="email"
              value={formData.email}
              onChange={handleChange}
            />
          </div>

          <div className="form-actions">
            <button
              type="submit"
              disabled={updateProfileMutation.isLoading}
              className="save-button"
            >
              {updateProfileMutation.isLoading ? 'Saving...' : 'Save Changes'}
            </button>
            <button
              type="button"
              onClick={() => setIsEditing(false)}
              className="cancel-button"
            >
              Cancel
            </button>
          </div>

          {updateProfileMutation.isError && (
            <p className="error-message">{updateProfileMutation.error?.message}</p>
          )}
        </form>
      )}
    </div>
  );
};

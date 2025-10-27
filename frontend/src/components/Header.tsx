import { useState } from 'react';
import { useCurrentUser } from '../contexts/CurrentUserContext';
import { useTransactionRefresh } from '../contexts/TransactionRefreshContext';
import { Select, SelectItem, Button } from '@tremor/react';
import { QuestionMarkCircleIcon } from '@heroicons/react/24/outline';
import { TransactionModal } from './TransactionModal';
import type { TransactionType } from '../types/transaction';

interface HeaderProps {
  onStartTour?: () => void;
}

export function Header({ onStartTour }: HeaderProps) {
  const { users, currentUser, setCurrentUserById, refreshUser, isLoading } = useCurrentUser();
  const { refreshTransactions } = useTransactionRefresh();
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [modalType, setModalType] = useState<TransactionType>('fund');

  if (isLoading) {
    return (
      <header className="bg-white dark:bg-gray-900 border-b border-gray-200 dark:border-gray-800 shadow-sm px-6 py-4">
        <div className="animate-pulse flex items-center space-x-4">
          <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-32"></div>
          <div className="h-8 bg-gray-200 dark:bg-gray-700 rounded w-48"></div>
          <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-40"></div>
        </div>
      </header>
    );
  }

  return (
    <header className="bg-white dark:bg-gray-900 border-b border-gray-200 dark:border-gray-800 shadow-sm px-4 sm:px-6 py-4">
      <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4 md:gap-6">
        {/* Logo/Title Section */}
        <div className="flex items-center gap-2">
          <h1 className="text-xl font-semibold text-gray-900 dark:text-gray-100">Dylan's Treasury Dashboard</h1>
        </div>

        {/* Right Section: User Dropdown, Balance, Actions */}
        <div className="flex flex-col sm:flex-row sm:items-center gap-3 sm:gap-4">
          {/* User Dropdown */}
          <div className="flex items-center gap-2" data-tour-id="user-selector">
            <span className="text-sm text-gray-600 dark:text-gray-400 whitespace-nowrap">User:</span>
            <Select
              value={currentUser?.id.toString() || ''}
              onValueChange={(value) => setCurrentUserById(parseInt(value, 10))}
              className="w-full sm:w-48"
            >
              {users.map((user) => (
                <SelectItem key={user.id} value={user.id.toString()}>
                  {user.name}
                </SelectItem>
              ))}
            </Select>
          </div>

          {/* Balance Display */}
          <div className="flex items-center gap-2 px-3 py-1.5 bg-gray-50 dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700" data-tour-id="balance-display">
            <span className="text-sm text-gray-600 dark:text-gray-400 whitespace-nowrap">Balance:</span>
            <div className="text-lg font-semibold text-gray-900 dark:text-gray-100">
              {currentUser
                ? `$${currentUser.balance.toLocaleString('en-US', {
                    minimumFractionDigits: 2,
                    maximumFractionDigits: 2,
                  })}`
                : '-'}
            </div>
          </div>

          {/* Action Buttons */}
          <div className="flex gap-2" data-tour-id="fund-withdraw-buttons">
            {/* Tour Button */}
            {onStartTour && (
              <button
                onClick={onStartTour}
                className="flex items-center justify-center p-2 text-gray-600 hover:text-gray-900 dark:text-gray-400 dark:hover:text-gray-100 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-lg transition-colors focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
                aria-label="View app walkthrough"
                title="View app walkthrough"
              >
                <QuestionMarkCircleIcon className="h-6 w-6" />
              </button>
            )}

            <Button
              onClick={() => {
                setModalType('fund');
                setIsModalOpen(true);
              }}
              variant="primary"
              className="flex-1 sm:flex-none"
            >
              Fund
            </Button>

            <Button
              onClick={() => {
                setModalType('withdraw');
                setIsModalOpen(true);
              }}
              variant="secondary"
              className="flex-1 sm:flex-none"
            >
              Withdraw
            </Button>
          </div>
        </div>
      </div>

      <TransactionModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        type={modalType}
        currentBalance={currentUser?.balance.toString() || '0'}
        userId={currentUser?.id || 0}
        onSuccess={async () => {
          await refreshUser();
          refreshTransactions();
        }}
      />
    </header>
  );
}

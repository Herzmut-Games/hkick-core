class GoalsController < ApplicationController
  def create
    case params[:team]
    when 'red'
      broadcast_goal('red')
    when 'white'
      broadcast_goal('white')
    end
  end

  private

  def broadcast_goal(team)
    ActionCable.server.broadcast 'goal_notifications_channel', team: team
  end
end
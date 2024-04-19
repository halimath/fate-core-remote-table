from playwright.sync_api import Page, Locator

class Scene:
    def __init__(self, page: Page, url: str | None = None):
        self._page = page
        self._url = url

    def visit_page(self):
        self._page.goto(self._url)

    def title(self) -> Locator:
        return self._page.locator('[data-testid="title"]')
    
    def skill_check_btn(self, skill_level: int) -> Locator:
        return self._page.locator(f'[data-testid="skill-check-{skill_level}-btn"]')
    
    def skill_check_result(self) -> Locator:
        return self._page.locator('[data-testid="skill-check-result"]')


class IndexScene(Scene):
    def __init__(self, page: Page):
        super(IndexScene, self).__init__(page, '/')

    def create_new_session_btn(self) -> Locator:
        return self._page.locator('[data-testid="create-session-btn"]')

    def join_session_btn(self) -> Locator:
        return self._page.locator('[data-testid="join-session-btn"]')


class GameMasterScene(Scene):
    class PlayerCard:
        def __init__(self, locator: Locator):
            self.locator = locator

        def fate_points_label(self) -> Locator:
            return self.locator.locator('[data-testid="fate-points"]')

        def dec_fate_points_btn(self) -> Locator:
            return self.locator.locator('[data-testid="dec-fate-points"]')

        def inc_fate_points_btn(self) -> Locator:
            return self.locator.locator('[data-testid="inc-fate-points"]')

    def join_session_link(self):
        return self._page.locator('[data-testid="join-session-link"]')

    def add_aspect_btn(self):
        return self._page.locator('[data-testid="add-aspect"]')

    def player_locator(self, i: int) -> Locator:
        return self._page.locator(f'[data-testid="players"] > div:nth-child({i+1})')

    def player(self, i: int) -> PlayerCard:
        return GameMasterScene.PlayerCard(self.player_locator(i))
    
    def aspect_locator(self, i: int) -> Locator:
        return self._page.locator(f'[data-testid="aspects"] > div:nth-child({i+1})')

    # def players(self) -> Sequence[PlayerCard]:
    #     return (GameMasterScene.PlayerCard(loc) for loc in self._page.locator('[data-testid="players"] > div').all())


class PlayerScene(Scene):
    def fate_points(self) -> Locator:
        return self._page.locator('[data-testid="fate-points"]')

    def spend_fate_points_btn(self) -> Locator:
        return self._page.locator('[data-testid="spend-fate-point"]')


class Modal:
    def __init__(self, page: Page):
        self.locator = page.locator('[data-testid="modal"]')

    def ok_btn(self) -> Locator:
        return self.locator.get_by_role('button', name='OK')


class CreateSessionModal(Modal):
    def session_title_input(self) -> Locator:
        return self.locator.locator('input[data-testid="session-title"]')


class AddAspectModal(Modal):
    def aspect_name_input(self) -> Locator:
        return self.locator.locator('input[data-testid="aspect-name"]')


class JoinSessionModal(Modal):
    def session_id_input(self) -> Locator:
        return self.locator.locator('input[data-testid="session-id"]')

    def player_name_input(self) -> Locator:
        return self.locator.locator('input[data-testid="player-name"]')
